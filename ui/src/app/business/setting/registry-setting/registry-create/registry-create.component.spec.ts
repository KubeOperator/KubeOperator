import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RegistryCreateComponent } from './registry-create.component';

describe('RegistryCreateComponent', () => {
  let component: RegistryCreateComponent;
  let fixture: ComponentFixture<RegistryCreateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RegistryCreateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RegistryCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
