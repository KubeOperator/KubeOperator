import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RegistryUpdateComponent } from './registry-update.component';

describe('RegistryUpdateComponent', () => {
  let component: RegistryUpdateComponent;
  let fixture: ComponentFixture<RegistryUpdateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RegistryUpdateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RegistryUpdateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
