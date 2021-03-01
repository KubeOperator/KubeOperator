import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RegistryDeleteComponent } from './registry-delete.component';

describe('RegistryDeleteComponent', () => {
  let component: RegistryDeleteComponent;
  let fixture: ComponentFixture<RegistryDeleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RegistryDeleteComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RegistryDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
