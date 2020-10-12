import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VmConfigCreateComponent } from './vm-config-create.component';

describe('VmConfigCreateComponent', () => {
  let component: VmConfigCreateComponent;
  let fixture: ComponentFixture<VmConfigCreateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ VmConfigCreateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(VmConfigCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
